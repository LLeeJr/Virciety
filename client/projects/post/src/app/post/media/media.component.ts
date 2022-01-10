import {Component, EventEmitter, Inject, Input, OnInit, Output} from '@angular/core';
import {Post} from "../../model/post";
import {MAT_DIALOG_DATA, MatDialog} from "@angular/material/dialog";

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

  constructor(private dialog: MatDialog) { }

  ngOnInit(): void {
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
  templateUrl: './dialog-liked-by.html',
})
export class DialogLikedBy {
  constructor(@Inject(MAT_DIALOG_DATA) public data: string[]) {}
}
