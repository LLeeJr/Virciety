import {Component, Inject, OnDestroy, OnInit} from '@angular/core';
import { Post } from "../model/post";
import {GQLService} from "../service/gql.service";
import {Observable} from "rxjs";
import {AuthLibService} from "auth-lib";
import {MAT_DIALOG_DATA, MatDialog} from "@angular/material/dialog";

@Component({
  selector: 'app-post',
  templateUrl: './post.component.html',
  styleUrls: ['./post.component.scss']
})
export class PostComponent implements OnInit, OnDestroy {
  posts: Observable<Post[]> | undefined;

  constructor(private gqlService: GQLService,
              private authService: AuthLibService,
              private dialog: MatDialog) {
  }

  ngOnInit(): void {
    this.posts = this.gqlService.getPosts();
  }

  ngOnDestroy(): void {
    this.gqlService.oldestPostReached = false;
  }

  get oldestPostReached(): boolean {
    return this.gqlService.oldestPostReached;
  }

  onScroll() {
    if (!this.oldestPostReached) {
      this.gqlService.refreshPosts();
    }
  }

  likePost(post: Post) {
    // TODO uncomment
    this.gqlService.likePost(post/*, this.authService.userName*/);
  }

  openLikedByDialog(likedBy: string[]) {
    this.dialog.open(DialogLikedBy, {
      data: likedBy
    });
  }

  editPost(id: string, newDescription: string) {
    this.gqlService.editPost(id, newDescription);
  }
}

@Component({
  selector: 'dialog-liked-by',
  templateUrl: './dialog-liked-by.html',
})
export class DialogLikedBy {
  constructor(@Inject(MAT_DIALOG_DATA) public data: string[]) {}
}
