import { Injectable } from '@angular/core';
import {Apollo, ApolloBase, gql, QueryRef} from "apollo-angular";
import {InMemoryCache} from "@apollo/client/core";
import {Post} from "../model/post";
import {DataLibService} from "data-lib";
import {HttpLink} from "apollo-angular/http";

@Injectable({
  providedIn: 'root'
})
export class GQLService {
// ------------------------- Queries, Mutations and Subscriptions

  private GET_POSTS = gql`
    query getPosts($id: String!, $fetchLimit: Int!) {
      getPosts(id: $id, fetchLimit: $fetchLimit) {
        id
        data {
          name
          contentType
        }
        description
        comments
        likedBy
      }
    }
  `;

  private GET_DATA = gql`
    query getData($id: String!) {
      getData(id: $id)
    }
  `;

  private CREATE_POST = gql`
    mutation createPost($username: String!, $description: String!, $data: String!) {
      createPost(newPost: {username: $username, description: $description, data: $data}) {
        id
        description
        data {
          name
          contentType
        }
        likedBy
        comments
      }
    }
  `;

  // ------------------------- Queries, Mutations and Subscriptions end

  private apollo: ApolloBase;
  private getPostQuery: QueryRef<any, {
    id: string;
    fetchLimit: number;
  }> | undefined;
  private lastPostID: string = "";
  private fetchLimit: number = 1;

  constructor(private apolloProvider: Apollo,
              private httpLink: HttpLink,
              private dataService: DataLibService) {
    const http = httpLink.create({
      uri: 'http://localhost:8083/query',
    });

    this.apolloProvider.createNamed('post', {
      cache: new InMemoryCache(),
      link: http,
    });

    this.apollo = this.apolloProvider.use('post');

    this.getPosts();
  }

  getPosts() {
    this.getPostQuery = this.apollo
      .watchQuery({
        query: this.GET_POSTS,
        variables: {
          id: this.lastPostID,
          fetchLimit: this.fetchLimit
        },
      });

    this.getPostQuery.valueChanges.subscribe(({data}: any) => {
      console.log(data);

      for (let getPost of data.getPosts) {
        const post: Post = new Post(getPost);

        this.lastPostID = post.id;
        this.dataService.addPost(post);
        this.getData(post);
      }
    }, (error: any) => {
      console.error('there was an error sending the getPost-query', error);
    });
  }

  refreshPosts() {
    if (this.getPostQuery) {
      this.getPostQuery.setVariables({
        id: this.lastPostID,
        fetchLimit: this.fetchLimit,
      }).then(() => {
        this.getPostQuery?.refetch()
      });
    }
  }

  getData(post: Post) {
    this.apollo.watchQuery({
      query: this.GET_DATA,
      variables: {
        id: post.id,
      },
    }).valueChanges.subscribe(({data}: any) => {
      post.data.content = data.getData;
    }, (error: any) => {
      console.error('there was an error sending the getData-query', error);
    });
  }

  async createPost(fileBase64: string, description: string, username: string) {
    this.apollo.mutate({
      mutation: this.CREATE_POST,
      variables: {
        username: username,
        description: description,
        data: fileBase64
      }
    }).subscribe((data: any) => {
      console.log('got data', data.data);
      const post = new Post(data.data.createPost);

      post.data.content = fileBase64;
      this.dataService.addNewPost(post);
    }, (error: any) => {
      console.error('there was an error sending the createPost-mutation', error);
    });
  }
}
