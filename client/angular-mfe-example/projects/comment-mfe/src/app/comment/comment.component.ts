import { Component, OnInit } from '@angular/core';
import {Apollo, gql} from "apollo-angular";
import { Comment } from "../model/comment";

@Component({
  selector: 'app-comment',
  templateUrl: './comment.component.html',
  styleUrls: ['./comment.component.css']
})
export class CommentComponent implements OnInit {

  loading = true;
  error: any;
  comments: Map<String, Comment[]> = new Map();

  constructor(private apollo: Apollo) { }

  ngOnInit(): void {
    this.apollo
      .watchQuery({
        query: gql`
          {
            GetComments {
              key
              value {
                id
                description
                likedBy
              }
            }
          }`,
      })
      .valueChanges.subscribe((data: any) => {
        console.log(data.data);

        let mapComments = data.data.GetComments

        for (let mapComment of mapComments) {
          this.comments.set(mapComment.key, mapComment.value)
        }

        this.loading = data.loading;
        this.error = data.error;
    })
  }

}
