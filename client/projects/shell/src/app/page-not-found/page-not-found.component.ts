import { Component, OnInit } from '@angular/core';
import {ActivatedRoute, Params} from "@angular/router";
import {MatSnackBar} from "@angular/material/snack-bar";

@Component({
  selector: 'app-page-not-found',
  templateUrl: './page-not-found.component.html',
  styleUrls: ['./page-not-found.component.scss']
})
export class PageNotFoundComponent implements OnInit {

  private durationTime: number = 3;

  constructor(private route: ActivatedRoute,
              private snackbar: MatSnackBar) { }

  ngOnInit(): void {
    this.route.params.subscribe((params: Params) => {
      let msg = params['msg'];
      this.snackbar.open(msg, undefined, {
        duration: this.durationTime * 1000,
      });
    });
  }

}
