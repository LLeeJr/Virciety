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
    GQLService._oldestPostReached = false;
  }

  get oldestPostReached(): boolean {
    // TODO when changing components and coming back to this one
    // check whether it is necessary to refetch newer posts than the newest this client has
    // ideas: send newest post id to server and ask whether new posts are there
    // subscription notification maybe?
    //
    return GQLService._oldestPostReached;
  }

  onScroll() {
    if (!this.oldestPostReached) {
      this.gqlService.refreshPosts();
    }
  }

  likePost(post: Post) {
    // TODO uncomment this
    const username: string = /*this.authService.userName*/ 'user4';
    let liked: boolean = true;

    // check if it's a like or unlike
    const index = post.likedBy.indexOf(username, 0);
    if (index > -1) {
      const newLikedBy: string[] = []
      post.likedBy.forEach((user, i) => {
        if (index !== i) {
          newLikedBy.push(user);
        }
      });
      post.likedBy = newLikedBy;
      liked = false;
    } else {
      const newLikedBy: string[] = [username];
      post.likedBy.forEach(user => newLikedBy.push(user));
      post.likedBy = newLikedBy;
    }

    this.gqlService.likePost(post, liked);
  }

  openLikedByDialog(likedBy: string[]) {
    this.dialog.open(DialogLikedBy, {
      data: likedBy
    });
  }

  editPost(post: Post) {
    this.gqlService.editPost(post);
  }

  removePost(post: Post) {
    this.gqlService.removePost(post);
  }
}

@Component({
  selector: 'dialog-liked-by',
  templateUrl: './dialog-liked-by.html',
})
export class DialogLikedBy {
  constructor(@Inject(MAT_DIALOG_DATA) public data: string[]) {}
}
