import {Component, EventEmitter, Inject, Input, OnInit, Output} from '@angular/core';
import {Post} from "../../model/post";
import {MAT_DIALOG_DATA, MatDialog} from "@angular/material/dialog";
import {AuthLibService} from "auth-lib";

@Component({
  selector: 'post-media',
  templateUrl: './media.component.html',
  styleUrls: ['./media.component.scss']
})
export class MediaComponent implements OnInit {

  @Input() post: Post;
  @Input() username: string;
  @Output() newEvent = new EventEmitter<string>();
  editMode: boolean = false;
  source: string = '';

  constructor(private auth: AuthLibService,
              private dialog: MatDialog) { }

  ngOnInit(): void {
    this.auth.getUserByName(this.post.username).subscribe(value => {
      if (value && value.data && value.data.getUserByName) {
        let {profilePictureId} = value.data.getUserByName;
        if (profilePictureId && profilePictureId !== '') {
          this.auth.getProfilePicture(profilePictureId).subscribe(picture => {
            if (picture && picture.data && picture.data.getProfilePicture) {
              this.source = picture.data.getProfilePicture;
            }
          });
        }
      }
    })
  }

  triggerEvent(eventName: string) {
    this.newEvent.emit(eventName);
  }

  openLikedByDialog(likedBy: string[]) {
    this.dialog.open(DialogLikedBy, {
      data: likedBy
    });
  }
}

@Component({
  selector: 'dialog-liked-by',
  template: `
    <h1 mat-dialog-title>Liked by</h1>
    <div mat-dialog-content>
        <mat-list>
            <mat-list-item role="listitem" *ngFor="let username of data">{{username}}</mat-list-item>
        </mat-list>
    </div>`,
})
export class DialogLikedBy {
  constructor(@Inject(MAT_DIALOG_DATA) public data: string[]) {}
}
