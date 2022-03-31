import {Component, EventEmitter, Inject, Input, OnInit, Output} from '@angular/core';
import {Post} from "../../model/post";
import {MAT_DIALOG_DATA, MatDialog} from "@angular/material/dialog";
import {AuthLibService} from "auth-lib";
import {SinglePostComponent} from "../single-post/single-post.component";
import { Location } from '@angular/common'
import {BreakpointObserver, Breakpoints} from "@angular/cdk/layout";

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
  isPhonePortrait: boolean = false;

  constructor(private auth: AuthLibService,
              private dialog: MatDialog,
              private location: Location,
              private responsive: BreakpointObserver) { }

  ngOnInit(): void {
    this.responsive.observe(Breakpoints.HandsetPortrait).subscribe((result) => {
      this.isPhonePortrait = result.matches;
    });

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

  openPostDialog(post: Post) {
    let state = this.location.path();
    this.location.replaceState(`/p/${post.id}`);
    let dialogRef = this.dialog.open(SinglePostComponent, {
      data: {
        post: post
      }
    })

    dialogRef.afterClosed().subscribe({
      next: _ => {
        this.location.replaceState(state);
      }
    })
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
