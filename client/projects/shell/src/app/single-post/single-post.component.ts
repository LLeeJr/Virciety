import { Component, OnInit } from '@angular/core';
import {environment} from "../../environments/environment";

@Component({
  selector: 'app-single-post',
  templateUrl: './single-post.component.html',
  styleUrls: ['./single-post.component.scss']
})
export class SinglePostComponent implements OnInit {
  postMFE: string;

  constructor() {
    this.postMFE = environment.postMFE;
  }

  ngOnInit(): void {
  }

}
