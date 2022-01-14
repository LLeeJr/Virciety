import {Component, OnInit} from '@angular/core';
import {KeycloakService} from "keycloak-angular";
import {FormControl} from "@angular/forms";
import {Observable} from "rxjs";
import {debounceTime, filter, map, startWith, tap} from "rxjs/operators";
import {AuthLibService} from "auth-lib";
import {Router} from "@angular/router";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'shell';
  isLoggedIn: boolean = false;
  username: string

  searchMode: boolean = false;
  searchFormControl = new FormControl();
  options: string[] = [];
  filteredOptions: Observable<string[]>;

  constructor(
    private auth: AuthLibService,
    private keycloak: KeycloakService,
    private router: Router,
  ) {
    this.keycloak.isLoggedIn().then((loggedIn) => {
      this.isLoggedIn = loggedIn;
      this.keycloak.loadUserProfile().then(() => this.username = this.keycloak.getUsername())
    });
  }

  ngOnInit() {}

  logout() {
    this.keycloak.isLoggedIn().then((loggedIn) => {
      if (loggedIn) {
        this.keycloak.logout(window.location.origin).then(() => {
          this.isLoggedIn = false;
        });
      }
    })
  }

  search(username: string) {
    this.auth.findUsersWithName(username).subscribe(value => {
      if (value && value.data && value.data.findUsersWithName) {
        let users: string[] = [];
        for (let user of value.data.findUsersWithName) {
          users.push(user.username);
        }
        this.options = [...users];
        this.filteredOptions = this.searchFormControl.valueChanges.pipe(
          startWith(''),
          map(value => this._filter(value)),
        );
      }
    }, () => {
      this.options = [];
    });
  }

  openSearch() {
    this.searchMode = !this.searchMode;
    this.searchFormControl.valueChanges.pipe(
      debounceTime(250),
    ).subscribe(username => {
      if (username && username.length > 0) {
        this.search(username);
      } else {
        this.options = [];
      }
    });
  }

  closeSearch() {
    this.searchMode = !this.searchMode;
    this.options = [];
    this.searchFormControl.reset();
  }

  private _filter(value: string) {
    if (this.options.length == 0) {
      return [];
    }
    const filterValue = value.toLowerCase();
    return this.options.filter(option => option.toLowerCase().includes(filterValue));
  }

  redirect(option: string) {
    if (option.length > 0) {
      this.router.navigate(['profile'], {
        queryParams: {
          username: option,
        }
      });
    }
  }
}
