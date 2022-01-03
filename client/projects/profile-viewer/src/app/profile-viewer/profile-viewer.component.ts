import {Component, EventEmitter, Inject, OnInit, Output} from '@angular/core';
import {AuthLibService} from "auth-lib";
import {ApiService, User} from "../api/api.service";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material/dialog";
import {Router} from "@angular/router";

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
              private auth: AuthLibService,
              private router: Router) { }

  ngOnInit(): void {
    console.log(this.router.url);
    this.router.navigate(['/home'])
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
    dialogRef.componentInstance.fileId.subscribe(fileId => {
      this.api.getUserByID(this.id).subscribe(value => {
        if (value && value.data && value.data.getUserByID) {
          this.activeUser = value.data.getUserByID;
        }
      });
      if (fileId !== '') {
        this.getProfilePicture(fileId);
      }
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
  @Output() fileId = new EventEmitter<string>();

  constructor(@Inject(MAT_DIALOG_DATA) public data: any,
              private api: ApiService,
              private dialogRef: MatDialogRef<ProfilePictureDialog>) {
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
            this.fileId.emit(value.data.addProfilePicture);
            this.dialogRef.close();
          }
        }, () => alert('Error during updating profile picture'));
    }
  }

  remove() {
    if (this.data && this.data.username) {
      if (this.data.profilePictureId) {
        this.api.removeProfilePicture(this.data.username, this.data.profilePictureId)
          .subscribe(value => {
            if (value && value.data && value.data.removeProfilePicture) {
              this.fileId.emit('');
              this.dialogRef.close();
            }
          }, () => alert('Error during removing profile picture'));
      } else {
        alert('No profile picture available to remove')
      }
    }
  }

  reset() {
    this.fileBase64 = null;
    this.content_type = '';
    this.filename = undefined;
  }
}
