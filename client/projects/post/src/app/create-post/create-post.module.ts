import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { CreatePostRoutingModule } from './create-post-routing.module';
import {MaterialModule} from "../material.module";
import {CreatePostComponent, DialogCreatePostComponent} from "./create-post.component";
import {FormsModule} from "@angular/forms";

const EXPORTS = [
  CreatePostComponent,
  DialogCreatePostComponent
]

@NgModule({
  declarations: [
    ...EXPORTS
  ],
  imports: [
    CommonModule,
    CreatePostRoutingModule,
    MaterialModule,
    FormsModule,
  ],
  exports: [...EXPORTS],
})
export class CreatePostModule {
  static exports = EXPORTS;
}
