import { Component, OnInit } from '@angular/core';
import {ApiService} from "../../../../user/src/app/api/api.service";
import {AuthLibService} from "auth-lib";

@Component({
  selector: 'app-profile-viewer',
  templateUrl: './profile-viewer.component.html',
  styleUrls: ['./profile-viewer.component.scss']
})
export class ProfileViewerComponent implements OnInit {

  id: string = '';
  activeUser: any;

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
  }

}
