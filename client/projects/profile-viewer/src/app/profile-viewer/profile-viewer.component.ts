import {Component, EventEmitter, Inject, OnInit, Output} from '@angular/core';
import {ApiService} from "../../../../user/src/app/api/api.service";
import {AuthLibService} from "auth-lib";
import {User} from "../api/api.service";
import {MAT_DIALOG_DATA, MatDialog} from "@angular/material/dialog";

@Component({
  selector: 'app-profile-viewer',
  templateUrl: './profile-viewer.component.html',
  styleUrls: ['./profile-viewer.component.scss']
})
export class ProfileViewerComponent implements OnInit {

  id: string = '';
  activeUser: User;
  source: string = '';

  constructor(public dialog: MatDialog,
              private api: ApiService,
              private auth: AuthLibService) { }

  ngOnInit(): void {
    this.auth._activeId.subscribe(id => {
      this.id = id;
      this.api.getUserByID(this.id).subscribe(value => {
        if (value && value.data && value.data.getUserByID) {
          this.activeUser = value.data.getUserByID;
          this.getProfilePicture(this.activeUser.profilePictureId);
        }
      });
    });
  }

  getProfilePicture(fileId: string) {
    if (this.activeUser && fileId !== '') {
      this.api.getProfilePicture(fileId).subscribe(value => {
        if (value && value.data && value.data.getProfilePicture) {
          this.source = value.data.getProfilePicture
        }
      });
    }
  }

  openDialog() {
    let dialogRef = this.dialog.open(ProfilePictureDialog, {
      data: this.activeUser,
    });
    dialogRef.componentInstance.sourceChanged.subscribe(result => {
      this.source = result;
    })
  }
}

@Component({
  selector: 'profile-picture-dialog',
  templateUrl: './profile-picture-dialog.html'
})
export class ProfilePictureDialog {
  fileBase64: any;
  content_type: string = '';
  filename: string | undefined;
  source: string;
  @Output() sourceChanged = new EventEmitter<string>();

  constructor(@Inject(MAT_DIALOG_DATA) public data: any,
              private api: ApiService) {
  }

  onFileSelected(event: any) {
    // get selected file
    const file = event.target.files[0] as File;
    this.filename = file.name;

    // get file data as base64 string
    if (file) {
      const reader = new FileReader();
      reader.readAsDataURL(file);

      reader.onload = () => {
        if (reader.result) {
          const base64 = reader.result;
          const data: string = base64.toString().split(";base64,")[0];

          this.content_type = data.split(":")[1];

          this.fileBase64 = base64;
        }
      }

      reader.onloadend = () => {
        this.upload();
      }
    }
  }

  upload() {
    if (this.fileBase64) {
      this.api.addProfilePicture(this.content_type, this.fileBase64, this.data.username)
        .subscribe((value) => {
          if (value && value.data && value.data.addProfilePicture) {
            this.sourceChanged.emit(this.fileBase64)
          }
        }, () => alert('Error during updating profile picture'));
    }
  }

  reset() {
    this.fileBase64 = null;
    this.content_type = '';
    this.filename = undefined;
  }
}
