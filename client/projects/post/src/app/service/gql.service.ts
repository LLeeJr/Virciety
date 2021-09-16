import { Injectable } from '@angular/core';
import {Apollo, ApolloBase, gql} from "apollo-angular";
import {HttpLink} from "apollo-angular/http";
import {InMemoryCache} from "@apollo/client/core";
import {Post} from "../model/post";
import {DataLibService} from "data-lib";

@Injectable({
  providedIn: 'root'
})
export class GQLService {
// ------------------------- Queries, Mutations and Subscriptions

  private GET_POSTS = gql`
    {
      getPosts {
        id
        data {
          id
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
          id
          contentType
        }
        likedBy
        comments
      }
    }
  `;

  // ------------------------- Queries, Mutations and Subscriptions end

  private apollo: ApolloBase;

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
    console.log("Get Posts");
    this.apollo
      .watchQuery({
        query: this.GET_POSTS,
      }).valueChanges.subscribe((data: any) => {
        let posts = [];

        for (let getPost of data.data.getPosts) {
          const post: Post = new Post(getPost);

          this.getData(post.id);

          posts.push(post);
        }

        this.dataService.setPosts(posts.reverse());
      });
  }

  getData(id: string) {
    this.apollo.watchQuery({
      query: this.GET_DATA,
      variables: {
        id: id,
      },
    }).valueChanges.subscribe(({data}: any) => {
      for (let post of this.dataService.getPosts()) {
        if (post.id === id) {
          post.data.content = data.getData;
          post.data.fileUrl = `data:${post.data.contentType};base64,${data.getData}`
        }
      }
    });
  }

  createPost(fileBase64: string, description: string, username: string) {
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

      post.data.fileUrl = fileBase64;
      this.dataService.addPost(post);
    }, (error: any) => {
      console.error('there was an error sending the query', error);
    });
  }
}
