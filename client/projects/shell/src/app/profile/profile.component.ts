import { Component, OnInit } from '@angular/core';
import {Router} from "@angular/router";
import {MatSnackBar} from "@angular/material/snack-bar";
import {environment} from "../../environments/environment";

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss']
})
export class ProfileComponent implements OnInit {
  profileMFE: string;
  postMFE: string;
  private durationTime: number = 3;

  constructor(private router: Router,
              private snackbar: MatSnackBar) {
    this.postMFE = environment.postMFE;
    this.profileMFE = environment.profileMFE;
  }

  ngOnInit(): void {
  }

  handleError(event: any) {
    let {error, component} = event;
    if (error) {
      let msg = `${component} is currently offline!`;
      switch (component) {
        case 'post':
          this.router.navigate(['/page-not-found', msg])
          break;
        case 'profileViewer':
          this.snackbar.open(msg, undefined, {duration: this.durationTime*1000});
          break;
      }
    }
  }
}
