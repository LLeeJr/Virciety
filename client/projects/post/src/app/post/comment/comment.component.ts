import {Component, Input, OnInit} from '@angular/core';
import {Post} from "../../model/post";

@Component({
  selector: 'app-post-comment',
  templateUrl: './comment.component.html',
  styleUrls: ['./comment.component.scss']
})
export class CommentComponent implements OnInit {

  @Input() post: Post;
  comments: any = [];
  comment: string = "";

  constructor() {
  }

  ngOnInit(): void {
    this.comments = [
      {
        comment: 'Amazing pic!',
        createdBy: 'Cubalatino',
        createdAt:  new Date().getTime()
      },
      {
        comment: 'Good job, Sir',
        createdBy: 'Faboss',
        createdAt:  new Date().getTime()
      },
      {
        comment: 'Random Comment innit',
        createdBy: 'RandomUser123',
        createdAt:  new Date().getTime()
      },
    ];

    this.comments.push(...this.post.comments);
  }

  addComment() {
    // TODO uncomment this
    const newComment = {
      comment: this.comment,
      createdBy: 'user3', //this.authService.username,
      createdAt: new Date().getTime()
    };

    this.post.comments = [newComment, ...this.post.comments];
  }

}
