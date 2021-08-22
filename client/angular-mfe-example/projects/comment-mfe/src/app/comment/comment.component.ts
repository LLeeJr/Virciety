import { Component, OnInit } from '@angular/core';
import {Apollo, ApolloBase, gql} from "apollo-angular";
import { Comment } from "../model/comment";
import {DataLibService} from "data-lib";

@Component({
  selector: 'app-comment',
  templateUrl: './comment.component.html',
  styleUrls: ['./comment.component.css']
})
export class CommentComponent implements OnInit {

  private apollo: ApolloBase;
  loading = true;
  error: any;
  comments: Map<string, Comment[]> = new Map();

  constructor(private apolloProvider: Apollo,
              private service: DataLibService) {
    this.apollo = apolloProvider.use('comment');
  }

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

        this.service.comments = mapComments;

        this.loading = data.loading;
        this.error = data.error;
    })
  }

}
