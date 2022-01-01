import {Injectable} from '@angular/core';
import {Apollo, ApolloBase, QueryRef} from "apollo-angular";
import {InMemoryCache, split} from "@apollo/client/core";
import {Post} from "../model/post";
import {Comment} from "../model/comment";
import {DataLibService} from "data-lib";
import {HttpLink} from "apollo-angular/http";
import {getMainDefinition} from "@apollo/client/utilities";
import {
  ADD_COMMENT,
  CREATE_POST,
  EDIT_POST,
  GET_DATA,
  GET_POST_COMMENTS,
  GET_POSTS,
  LIKE_POST,
  NEW_POST_CREATED,
  REMOVE_POST
} from "./gql-request-strings";
import {WebSocketLink} from "@apollo/client/link/ws";
import {map} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class GQLService {
  private apollo: ApolloBase;
  private getPostQuery: QueryRef<any, {
    id: string;
    fetchLimit: number;
  }> | undefined;

  private lastPostID: string = "";
  private _fetchLimit: number = 5;
  static _oldestPostReached: boolean = false;

  constructor(private apolloProvider: Apollo,
              private httpLink: HttpLink,
              private dataService: DataLibService) {
    const http = httpLink.create({
      uri: 'http://localhost:8083/query',
    });

    const ws = new WebSocketLink({
      uri: `ws://localhost:8083/query`,
    })

    const link = split(
      ({query}) => {
        const data = getMainDefinition(query);
        return (
          data.kind === 'OperationDefinition' && data.operation === 'subscription'
        );
      },
      ws,
      http
    )

    try {
      this.apolloProvider.createNamed('post', {
        cache: new InMemoryCache({
          typePolicies: {
            Query: {
              fields: {
                getPosts: {
                  keyArgs: [],
                  merge(existing = [], incoming, { args}) {
                    // console.log('Existing: ', existing);
                    // console.log('Incoming: ', incoming);
                    // console.log('Args: ', args);

                    if (incoming.length === 0) {
                      GQLService._oldestPostReached = true;
                      return existing;
                    }

                    if (args && (args['id'] === 'remove' || args['id'] === 'create')) {
                      return incoming
                    }

                    return [...existing, ...incoming];
                  },
                },
              }
            }
          }
        }),
        link: link,
      });
    } catch (e) {
      console.error('Error when creating apollo client \'post\'', e);
    }

    this.apollo = this.apolloProvider.use('post');
  }

  // Getter + Setter

  get fetchLimit(): number {
    return this._fetchLimit;
  }

  // Getter + Setter end

  getPosts(): any {
    if (!this.getPostQuery) {
      this.getPostQuery = this.apollo
        .watchQuery({
          query: GET_POSTS,
          variables: {
            id: this.lastPostID,
            fetchLimit: this.fetchLimit
          },
        });
    }

    return this.getPostQuery?.valueChanges.pipe(map(({data}: any) => {
      // console.log('GetPostData: ', data);

      const posts = this.dataService.posts;

      for (let getPost of data.getPosts) {
        if (posts.some(p => p.id === getPost.id)) {
          // console.log('post already exists');
          continue;
        }

        const post: Post = new Post(getPost);
        this.lastPostID = post.id;

        this.getData(post);
        posts.push(post);
      }

      return posts;
    }, (error: any) => {
      console.error('there was an error sending the getPost-query', error);
    }));
  }

  refreshPosts() {
    if (this.getPostQuery) {
      this.getPostQuery.fetchMore({
        variables: {
          id: this.lastPostID,
          fetchLimit: this.fetchLimit,
        }
      }).then(() => {
        // console.log('getPost fetchMore success!');
      },
        error => {
        console.error('there was an error refreshing getPost-query', error);
        });
    } else {
      console.error('getPostQuery is null|undefined');
    }
  }

  private async getData(post: Post) {
    this.apollo.watchQuery({
      query: GET_DATA,
      variables: {
        fileID: post.data.name,
      },
    }).valueChanges.subscribe(({data}: any) => {
      post.data.content = data.getData;
    }, (error: any) => {
      console.error('there was an error sending the getData-query', error);
    });
  }

  getPostComments(post: Post) {
    this.apollo.watchQuery({
      query: GET_POST_COMMENTS,
      fetchPolicy: "network-only",
      variables: {
        id: post.id
      },
    }).valueChanges.subscribe(({data}: any) => {
      // console.log('GetPostComments-data: ', data);

      const commentList: Comment[] = [];

      for (let getPostComment of data.getPostComments) {
        commentList.push(new Comment(getPostComment));
      }

      post.comments = commentList;

      post.commentMode = true;
    }, (error: any) => {
      console.error('there was an error sending the getPostComments-query', error);
    })
  }

  // when a post is created before posts are fetched from server, smth doesn't work as it should
  // could be not a problem in usage with mfe
  async createPost(contentType: string, fileBase64: string, description: string, username: string = 'user3') {
    this.apollo.mutate({
      mutation: CREATE_POST,
      variables: {
        username: username,
        description: description,
        data: fileBase64
      },
      update: (cache, {data}: any) => {
        const existingPosts: any = cache.readQuery({
          query: GET_POSTS
        });

        const post = new Post(data.createPost);

        if (this.dataService.posts.some(p => p.id === post.id)) {
          // console.log('post already exists');
          return;
        }

        post.data.content = fileBase64;
        this.dataService.addNewPost(post);

        let newPosts;

        if (existingPosts) {
          newPosts = [data.createPost, ...existingPosts.getPosts];
        } else {
          newPosts = [data.createPost];
        }

        cache.writeQuery({
          query: GET_POSTS,
          variables: {
            id: 'create',
          },
          data: { getPosts : newPosts }
        });
      },
      /*context: {
        useMultipart: true,
      }*/
      }).subscribe(({data}: any) => {
        // console.log('CreatePostData: ', data);
    }, (error: any) => {
      console.error('there was an error sending the createPost-mutation', error);
    });
  }

  likePost(post: Post, liked: boolean) {
    this.apollo.mutate({
      mutation: LIKE_POST,
      variables: {
        id: post.id,
        description: post.description,
        newLikedBy: post.likedBy,
        comments: post.comments,
        liked: liked,
      }
    }).subscribe(({data}: any) => {
      // console.log('LikePostData: ', data);
    }, (error: any) => {
      console.error('there was an error sending the likePost-mutation', error);
    })
  }

  editPost(post: Post) {
    this.apollo.mutate({
      mutation: EDIT_POST,
      variables: {
        id: post.id,
        newDescription: post.description,
        likedBy: post.likedBy,
        comments: post.comments,
      }
    }).subscribe(({data}) => {
      // console.log('EditPostData: ', data)
    }, (error: any) => {
      console.error('there was an error sending the editPost-mutation', error);
    });
  }

  removePost(post: Post) {
    this.apollo.mutate({
      mutation: REMOVE_POST,
      variables: {
        id: post.id,
        fileID: post.data.name
      },
      update: (cache) => {
        const existingPosts: any = cache.readQuery({
          query: GET_POSTS,
        });
        const newPosts = existingPosts.getPosts.filter((getPost: any) => getPost.id !== post.id);

        this.dataService.removePost(post.id);

        cache.writeQuery({
          query: GET_POSTS,
          variables: {
            id: 'remove'
          },
          data: { getPosts : newPosts }
        });
      },
    }).subscribe(({data}: any) => {
      // console.log('RemovePostData: ', data);
    }, (error: any) => {
      console.error('there was an error sending the removePost-mutation', error);
    });
  }

  getPostCreated() {
    this.apollo.subscribe({
      query: NEW_POST_CREATED,
    }).subscribe(({data}: any) => {
      // console.log('NewPostCreated: ', data);

      const post = new Post(data.newPostCreated);

      if (this.dataService.posts.some(p => p.id === post.id)) {
        // console.log('post already exists');
        return;
      }

      this.dataService.addNewPost(post);
      this.getData(post);

      const cache = this.apollo.client.cache;

      const existingPosts: any = cache.readQuery({
        query: GET_POSTS,
      });

      let newPosts;

      if (existingPosts) {
        newPosts = [data.newPostCreated, ...existingPosts.getPosts];
      } else {
        newPosts = [data.newPostCreated];
      }

      cache.writeQuery({
        query: GET_POSTS,
        variables: {
          id: 'create',
        },
        data: { getPosts : newPosts }
      });
    }, (error: any) => {
      console.error('there was an error sending the newPostCreated-subscription', error)
    })
  }

  addComment(post: Post, addCommentRequest: { createdBy: string; comment: string; postID: string }) {
    this.apollo.mutate({
      mutation: ADD_COMMENT,
      variables: {
        comment: addCommentRequest,
      }
    }).subscribe(({data}: any) => {
      const comment = new Comment(data.addComment);

      post.comments = [comment, ...post.comments];

      // console.log('AddCommentData: ', data)
    }, (error: any) => {
      console.error('there was an error sending the addComment-mutation', error);
    });
  }
}
