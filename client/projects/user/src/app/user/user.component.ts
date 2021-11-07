import { Component, OnInit } from '@angular/core';
import {ApiService} from "../api/api.service";

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.scss']
})
export class UserComponent implements OnInit {

  constructor(private api: ApiService) { }

  ngOnInit(): void {
    // this.api.getUserByName("fabeeey").subscribe(value => console.log(value));
    // this.api.addFollow("618054cb5eaeb32e41c60530", "bob").subscribe(value => console.log(value));
    // this.api.getUserByID("618054cb5eaeb32e41c60530").subscribe(value => console.log(value));
    // this.api.removeFollow("618054cb5eaeb32e41c60530", "bob").subscribe(value => console.log(value));
    // this.api.getUserByID("618054cb5eaeb32e41c60530").subscribe(value => console.log(value));
  }

}
