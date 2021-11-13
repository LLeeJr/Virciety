import { Component, OnInit } from '@angular/core';
import {ApiService} from "../api/api.service";
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

  constructor(private api: ApiService,
              private auth: AuthLibService) { }

  ngOnInit(): void {
    this.auth._activeId.subscribe(value => this.id = value);
    this.searchFormControl
      .valueChanges
      .pipe(
        filter((username) => username),
        debounceTime(250),
      )
      .subscribe(username => {
        this.search(username);
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
}
