import {Component, OnDestroy, OnInit} from '@angular/core';
import {AuthLibService, User} from "auth-lib";
import {FormControl} from "@angular/forms";
import {debounceTime, map, startWith} from "rxjs/operators";
import {Observable} from "rxjs";
import {Router} from "@angular/router";


@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.scss'],
  exportAs: 'UserComponent'
})
export class UserComponent implements OnInit, OnDestroy {

  searchFormControl = new FormControl();
  options: string[] = [];
  filteredOptions: Observable<string[]>;
  id: string = '';
  activeUser: User;

  constructor(private auth: AuthLibService,
              private router: Router) { }

  ngOnInit(): void {
    this.auth._activeId.subscribe(id => {
      this.id = id;
      this.auth.getUserByID(this.id).subscribe(value => {
        if (value && value.data && value.data.getUserByID) {
          this.activeUser = value.data.getUserByID;
        }
      });
    });
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

  ngOnDestroy(): void {
    this.options = [];
    this.searchFormControl.reset();
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
