import { Component, OnInit } from '@angular/core';
import {ApiService} from "../api/api.service";
import {AuthLibService} from "auth-lib";

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.scss']
})
export class UserComponent implements OnInit {

  id: string = '';

  constructor(private api: ApiService,
              private auth: AuthLibService) { }

  ngOnInit(): void {
    // this.api.getUserByName("fabeeey").subscribe(value => console.log(value));
    // this.api.addFollow("618054cb5eaeb32e41c60530", "bob").subscribe(value => console.log(value));
    // this.api.getUserByID("618054cb5eaeb32e41c60530").subscribe(value => console.log(value));
    // this.api.removeFollow("618054cb5eaeb32e41c60530", "bob").subscribe(value => console.log(value));
    // this.api.getUserByID("618054cb5eaeb32e41c60530").subscribe(value => console.log(value));
    this.auth._activeId.subscribe(value => this.id = value)
    // this.id = this.auth.activeId;
  }

}
