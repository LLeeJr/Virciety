import {Component, Inject, Input, OnInit} from '@angular/core';
import {Post} from "../../model/post";
import {MAT_DIALOG_DATA, MatDialog} from "@angular/material/dialog";
import {GQLService} from "../../service/gql.service";
import {AuthLibService} from "auth-lib";

@Component({
  selector: 'app-post-media',
  templateUrl: './media.component.html',
  styleUrls: ['./media.component.scss']
})
export class MediaComponent implements OnInit {

  @Input() post: Post
  editMode: boolean = false;

  constructor(private dialog: MatDialog,
              private gqlService: GQLService,
              private authService: AuthLibService) { }

  ngOnInit(): void {
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

  editPost(post: Post) {
    this.gqlService.editPost(post);
  }

  removePost(post: Post) {
    this.gqlService.removePost(post);
  }

  openLikedByDialog(likedBy: string[]) {
    this.dialog.open(DialogLikedBy, {
      data: likedBy
    });
  }

}

@Component({
  selector: 'dialog-liked-by',
  templateUrl: './dialog-liked-by.html',
})
export class DialogLikedBy {
  constructor(@Inject(MAT_DIALOG_DATA) public data: string[]) {}
}
