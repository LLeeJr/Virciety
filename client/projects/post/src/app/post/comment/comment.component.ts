import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Post} from "../../model/post";

@Component({
  selector: 'app-post-comment',
  templateUrl: './comment.component.html',
  styleUrls: ['./comment.component.scss']
})
export class CommentComponent implements OnInit {

  @Input() post: Post;
  @Output() newCommentEvent = new EventEmitter<string>()
  comment: string = "";

  constructor() {
  }

  ngOnInit(): void {
  }

  addComment() {
    this.newCommentEvent.emit(this.comment);
    this.comment = "";
  }

}
