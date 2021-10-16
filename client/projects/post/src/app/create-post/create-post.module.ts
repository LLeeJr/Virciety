import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { CreatePostRoutingModule } from './create-post-routing.module';
import {MaterialModule} from "../material.module";
import {CreatePostComponent} from "./create-post.component";
import {FormsModule} from "@angular/forms";


@NgModule({
    declarations: [
        CreatePostComponent
    ],
    imports: [
        CommonModule,
        CreatePostRoutingModule,
        MaterialModule,
        FormsModule,
    ]
})
export class CreatePostModule { }
