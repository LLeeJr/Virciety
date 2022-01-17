import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { EventRoutingModule } from './event-routing.module';
import { EventComponent } from './event.component';
import {MaterialModule} from "../material.module";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";


@NgModule({
  declarations: [
    EventComponent
  ],
    imports: [
        CommonModule,
        EventRoutingModule,
        MaterialModule,
        FormsModule,
        ReactiveFormsModule,
    ]
})
export class EventModule { }
