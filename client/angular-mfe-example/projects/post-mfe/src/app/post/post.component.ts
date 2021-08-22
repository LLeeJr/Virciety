import { Component, OnInit } from '@angular/core';
import {Apollo, ApolloBase, gql} from "apollo-angular";
import { Post } from "../model/post";
import {DataLibService} from "data-lib";

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

  constructor(private apolloProvider: Apollo,
              private service: DataLibService) {
    this.apollo = this.apolloProvider.use('post');
  }

  ngOnInit(): void {
    this.apollo
      .watchQuery({
        query: gql`
          {
            getPosts {
              id
              data
              description
              comments
              likedBy
            }
          }`,
      })
      .valueChanges.subscribe((data: any) => {
      this.service.posts = data.data.getPosts;
      this.posts = data.data.getPosts;
      this.loading = data.loading;
      this.error = data.error;
    });
  }

}
