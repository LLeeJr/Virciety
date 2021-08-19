import { Component, OnInit } from '@angular/core';
import { Apollo, gql } from "apollo-angular";
import { Post } from "../model/post";

@Component({
  selector: 'app-post',
  templateUrl: './post.component.html',
  styleUrls: ['./post.component.css']
})
export class PostComponent implements OnInit {

  loading = true;
  error: any;
  posts: Post[] = [];

  constructor(private apollo: Apollo) {
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
      this.posts = data.data.getPosts;
      this.loading = data.loading;
      this.error = data.error;
    });
  }

}
