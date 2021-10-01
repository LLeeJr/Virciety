import { Injectable } from '@angular/core';
import {Apollo, ApolloBase, QueryRef} from "apollo-angular";
import {InMemoryCache, split} from "@apollo/client/core";
import {Post} from "../model/post";
import {DataLibService} from "data-lib";
import {HttpLink} from "apollo-angular/http";
import {getMainDefinition} from "@apollo/client/utilities";
import {CREATE_POST, EDIT_POST, GET_DATA, GET_POSTS, LIKE_POST} from "./gql-request-strings";
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
  private _oldestPostReached: boolean = false;

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

    this.apolloProvider.createNamed('post', {
      cache: new InMemoryCache(),
      link: link,
    });

    this.apollo = this.apolloProvider.use('post');

    this.getPostCreated();
  }

  // Getter + Setter

  get oldestPostReached(): boolean {
    return this._oldestPostReached;
  }

  set oldestPostReached(value: boolean) {
    this._oldestPostReached = value;
  }

  get fetchLimit(): number {
    return this._fetchLimit;
  }

  // Getter + Setter end

  getPosts(): any {
    if (this.dataService.posts.length === 0) {
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
      console.log('GetPostData: ', data);

      if (data.getPosts.length == 0) {
        this._oldestPostReached = true;
      }

      const posts = this.dataService.posts;

      for (let getPost of data.getPosts) {
        if (posts.some(p => p.id === getPost.id)) {
          console.log("post already exists");
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
      this.getPostQuery.setVariables({
        id: this.lastPostID,
        fetchLimit: this.fetchLimit,
      }).then(() => {
        this.getPostQuery?.refetch();
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
        fileID: post.data.id,
      },
    }).valueChanges.subscribe(({data}: any) => {
      post.data.content = data.getData;
    }, (error: any) => {
      console.error('there was an error sending the getData-query', error);
    });
  }

  async createPost(fileBase64: string, description: string, username: string = 'user3') {
    this.apollo.mutate({
      mutation: CREATE_POST,
      variables: {
        username: username,
        description: description,
        data: fileBase64
      },
      context: {
        useMultipart: true,
      }
      }).subscribe(({data}: any) => {
      console.log('CreatePostData: ', data);
      const post = new Post(data.createPost);

      post.data.content = fileBase64;
      this.dataService.addNewPost(post);
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
      console.log('LikePostData: ', data);
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
      console.log('EditPostData: ', data)
    }, (error: any) => {
      console.error('there was an error sending the likePost-mutation', error);
    });
  }

  private getPostCreated() {
    /*this.getPostQuery?.subscribeToMore({
      document: NEW_POST_CREATED,
      updateQuery: (prev: any, {data}: any) => {
        console.log('GetPostCreatedData: ', data);
        console.log('GetPostCreatedData prev: ', prev);
      }
    })*/


    /*this.apollo.subscribe({
      query: NEW_POST_CREATED,
    }).subscribe(({data}: any) => {

      if (data.newPostCreated.username === 'user3') {
        return
      }
      const post = new Post(data.newPostCreated);

      this.getData(post);
    }, (error: any) => {
      console.error('there was an error sending the newPostCreated-subscription', error)
    });*/
  }
}
