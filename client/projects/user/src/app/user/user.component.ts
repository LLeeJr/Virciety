import { Component, OnInit } from '@angular/core';
import {ApiService, User} from "../api/api.service";
import {AuthLibService} from "auth-lib";
import {FormControl} from "@angular/forms";
import {debounceTime, filter} from "rxjs/operators";

class UserData {
  constructor(firstName: string, lastName: string, username: string) {
    this.firstName = firstName;
    this.lastName = lastName;
    this.username = username;
  }

  firstName: string;
  lastName: string;
  username: string;
}

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.scss']
})
export class UserComponent implements OnInit {

  id: string = '';
  searchFormControl = new FormControl();
  users: UserData[] = [];
  activeUser: User;

  constructor(private api: ApiService,
              private auth: AuthLibService) { }

  ngOnInit(): void {
    this.auth._activeId.subscribe(id => {
      this.id = id;
      this.api.getUserByID(this.id).subscribe(value => {
        if (value && value.data && value.data.getUserByID) {
          this.activeUser = value.data.getUserByID;
        }
      });
    });
    this.searchFormControl
      .valueChanges
      .pipe(
        filter((username) => username),
        debounceTime(250),
      )
      .subscribe(username => {
        if (username.length > 0) {
          this.search(username);
        } else {
          this.users = [];
        }
      }, () => {
        this.users = [];
      });
  }

  search(username: string) {
    if (username && username.length > 0) {
      this.api.findUsersWithName(username).subscribe(value => {
        if (value && value.data && value.data.findUsersWithName) {
          this.users = value.data.findUsersWithName;
        }
      });
    }
  }

  checkFollow(username: string): boolean {
    if (this.activeUser) {
      return this.activeUser.follows.includes(username);
    } else {
      return false;
    }
  }

  removeFollow(username: string) {
    let id = this.activeUser.id;
    if (id && username) {
      this.api.removeFollow(id, username).subscribe(value => {
        if (value && value.data && value.data.removeFollow) {
          this.activeUser = value.data.removeFollow;
        }
      });
    }
  }

  addFollow(username: string) {
    let id = this.activeUser.id;
    if (id && username) {
      this.api.addFollow(id, username).subscribe(value => {
        if (value && value.data && value.data.addFollow) {
          this.activeUser = value.data.addFollow;
          console.log(this.activeUser);
        }
      });
    }
  }

  followButtonText(username: string): string {
    return this.activeUser.follows.includes(username) ? "Unfollow" : "Follow";
  }
}
