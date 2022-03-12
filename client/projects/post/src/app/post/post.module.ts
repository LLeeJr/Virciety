import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { PostRoutingModule } from './post-routing.module';
import {PostComponent} from './post.component';
import {InfiniteScrollModule} from "ngx-infinite-scroll";
import {MaterialModule} from "../material.module";
import {FormsModule} from "@angular/forms";
import {MediaComponent, DialogLikedBy} from './media/media.component';
import { CommentComponent } from './comment/comment.component';
import { SinglePostComponent } from './single-post/single-post.component';

const EXPORTS = [
  PostComponent,
  SinglePostComponent
]

@NgModule({
  declarations: [
    ...EXPORTS,
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
  ],
  exports: [...EXPORTS],

})
export class PostModule {
  static exports = EXPORTS;
}
