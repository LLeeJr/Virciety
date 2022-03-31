import {AfterViewInit, Component, ElementRef, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {AuthLibService, User} from "auth-lib";
import {FormControl} from "@angular/forms";
import {debounceTime, map, startWith} from "rxjs/operators";
import {Observable} from "rxjs";
import {Router} from "@angular/router";
import {BreakpointObserver, Breakpoints} from "@angular/cdk/layout";
import {MatSnackBar} from "@angular/material/snack-bar";
import {MatAutocomplete} from "@angular/material/autocomplete";


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
  isPhonePortrait: boolean = false;
  searchMode: boolean = false;
  private reset: boolean = false;

  constructor(private auth: AuthLibService,
              private router: Router,
              private responsive: BreakpointObserver,
              private snackbar: MatSnackBar) { }

  ngOnInit(): void {
    if (!this.auth.error) {
      this.auth.errorState.subscribe(value => this.snackbar.open(value, undefined, {duration: 3000}));
    } else {
      this.snackbar.open(this.auth.error, undefined, {duration: 3000})
    }
    this.responsive.observe(Breakpoints.HandsetPortrait).subscribe((result) => {
      this.isPhonePortrait = result.matches;
    });

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

  handleBlur() {
    setTimeout(() => {
      if (this.reset) {
        this.reset = false;
      } else {
        this.searchMode = !this.searchMode;
      }
    }, 200);
  }

  handleReset() {
    this.reset = true;
  }
}
