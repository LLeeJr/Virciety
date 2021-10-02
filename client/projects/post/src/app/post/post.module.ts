import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { PostRoutingModule } from './post-routing.module';
import {DialogLikedBy, PostComponent} from './post.component';
import {InfiniteScrollModule} from "ngx-infinite-scroll";
import {MaterialModule} from "../material.module";
import {FormsModule} from "@angular/forms";


@NgModule({
  declarations: [
    PostComponent,
    DialogLikedBy
  ],
  imports: [
    CommonModule,
    PostRoutingModule,
    InfiniteScrollModule,
    MaterialModule,
    FormsModule
  ]
})
export class PostModule { }
