import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { PostRoutingModule } from './post-routing.module';
import {PostComponent} from './post.component';
import {InfiniteScrollModule} from "ngx-infinite-scroll";
import {MaterialModule} from "../material.module";
import {FormsModule} from "@angular/forms";
import {MediaComponent, DialogLikedBy} from './media/media.component';
import { CommentComponent } from './comment/comment.component';


@NgModule({
  declarations: [
    PostComponent,
    MediaComponent,
    DialogLikedBy,
    CommentComponent
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