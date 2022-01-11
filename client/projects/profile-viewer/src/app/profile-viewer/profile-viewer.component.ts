import {Component, EventEmitter, Inject, OnInit, Output} from '@angular/core';
import {AuthLibService, User} from "auth-lib";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material/dialog";
import {ActivatedRoute, Router} from "@angular/router";
import {KeycloakService} from "keycloak-angular";

@Component({
  selector: 'app-profile-viewer',
  templateUrl: './profile-viewer.component.html',
  styleUrls: ['./profile-viewer.component.scss'],
  exportAs: 'ProfileViewerComponent'
})
export class ProfileViewerComponent implements OnInit {

  id: string = '';
  activeUser: User;
  source: string = '';
  loggedInUser: string;
  pickedUser: string;

  constructor(public dialog: MatDialog,
              private auth: AuthLibService,
              private keycloak: KeycloakService,
              private route: ActivatedRoute,
              private router: Router) { }

  async ngOnInit(): Promise<void> {
    await this.keycloak.isLoggedIn().then(() => {
      this.keycloak.loadUserProfile().then(() => {
        this.loggedInUser = this.keycloak.getUsername();
        this.route.queryParams.subscribe(({username}) => {
          this.pickedUser = username;
          this.auth.getUserByName(username).subscribe(value => {
            if (value && value.data && value.data.getUserByName) {
              this.activeUser = value.data.getUserByName;
              this.getProfilePicture(this.activeUser.profilePictureId);
            }
          });
        })
        this.auth._activeId.subscribe(id => {
          this.id = id;
        });
      });
    });
  }

  getProfilePicture(fileId: string) {
    if (this.activeUser && fileId !== '') {
      this.auth.getProfilePicture(fileId).subscribe(value => {
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
      this.auth.getUserByID(this.activeUser.id).subscribe(value => {
        if (value && value.data && value.data.getUserByID) {
          this.activeUser = value.data.getUserByID;
        }
      }, (error => console.error(error)));
      if (fileId !== '') {
        this.getProfilePicture(fileId);
      }
    })
  }

  openChat() {
    this.router.navigate(['/chat', this.pickedUser]);
  }

  isCurrentUser(): boolean {
    return this.loggedInUser === this.pickedUser;
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
              private auth: AuthLibService,
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
      this.auth.addProfilePicture(this.content_type, this.fileBase64, this.data.username)
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
        this.auth.removeProfilePicture(this.data.username, this.data.profilePictureId)
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
