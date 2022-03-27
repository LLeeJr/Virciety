import {ChangeDetectorRef, Component, OnInit} from '@angular/core';
import {GQLService} from "../service/gql.service";
import {MatDialog} from "@angular/material/dialog";
import {KeycloakService} from "keycloak-angular";
import {BreakpointObserver, Breakpoints} from "@angular/cdk/layout";

@Component({
  selector: 'app-create-post',
  templateUrl: './create-post.component.html',
  styleUrls: ['./create-post.component.scss'],
  exportAs: 'CreatePostComponent'
})
export class CreatePostComponent implements OnInit {
  fileBase64: any;
  description: string = '';
  content_type: string = '';
  filename: string | undefined;
  username: string;
  isPhonePortrait: boolean = false;

  constructor(private gqlService: GQLService,
              private keycloak: KeycloakService,
              private cd: ChangeDetectorRef,
              private responsive: BreakpointObserver) {
  }

  async ngOnInit(): Promise<void> {
    this.responsive.observe(Breakpoints.HandsetPortrait).subscribe((result) => {
      this.isPhonePortrait = result.matches;
    });

    await this.keycloak.isLoggedIn().then(loggedIn => {
      if (loggedIn) {
        this.keycloak.loadUserProfile().then(() => {
          this.username = this.keycloak.getUsername();
        })
      } else {
        this.keycloak.login();
      }
    });
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
    }
  }

  createPost() {
    if (this.fileBase64) {
      this.gqlService.createPost(this.content_type, this.fileBase64, this.description, this.username).then(() => {
        this.reset();
        console.log("File upload complete");
      });
    }
  }

  alertFunction() {
    alert(`Content-Type ${this.content_type} is not supported`);
    this.content_type = '';
    this.cd.detectChanges();
  }

  reset() {
    this.fileBase64 = null;
    this.description = '';
    this.content_type = '';
    this.filename = undefined;
  }
}

@Component({
  selector: 'app-dialog-create-post',
  template: `
    <button mat-icon-button (click)="openDialog()">
      <mat-icon>add</mat-icon>
    </button>`,
  styleUrls: ['./create-post.component.scss'],
  exportAs: 'DialogCreatePostComponent'
})
export class DialogCreatePostComponent implements OnInit {
  constructor(private dialog: MatDialog) {
  }

  ngOnInit(): void {
  }

  openDialog() {
    this.dialog.open(CreatePostComponent);
  }
}

