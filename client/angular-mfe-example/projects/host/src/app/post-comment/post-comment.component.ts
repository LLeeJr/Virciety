import { Component, OnInit } from '@angular/core';
import {DataLibService} from "data-lib";
import {Post} from "../../../../data-lib/src/model/post";

@Component({
  selector: 'app-post-comment',
  templateUrl: './post-comment.component.html',
  styleUrls: ['./post-comment.component.css']
})
export class PostCommentComponent implements OnInit {

  posts: Post[] = [];

  constructor(private service: DataLibService) { }

  ngOnInit(): void {
    this.posts = this.service.posts;
  }

}
