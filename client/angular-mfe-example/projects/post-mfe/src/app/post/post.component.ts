import { Component, OnInit } from '@angular/core';
import {Apollo, ApolloBase, gql} from "apollo-angular";
import { Post } from "../model/post";

@Component({
  selector: 'app-post',
  templateUrl: './post.component.html',
  styleUrls: ['./post.component.css']
})
export class PostComponent implements OnInit {
  private apollo: ApolloBase;
  loading = true;
  error: any;
  posts: Post[] = [];

  constructor(private apolloProvider: Apollo) {
    this.apollo = this.apolloProvider.use('post');
  }

  ngOnInit(): void {
    this.apollo
      .watchQuery({
        query: gql`
          {
            getPosts {
              id
              data {
                id
                content
                contentType
              }
              description
              comments
              likedBy
            }
          }`,
      })
      .valueChanges.subscribe((data: any) => {

      console.log(data.data.getPosts);

      for (let getPost of data.data.getPosts) {
        const post: Post = new Post(getPost);
        this.posts.push(post)
      }

      this.loading = data.loading;
      this.error = data.error;
    });
  }

  createPost() {

  }
}
