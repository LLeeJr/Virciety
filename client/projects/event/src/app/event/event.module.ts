import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { EventRoutingModule } from './event-routing.module';
import { EventComponent } from './event.component';
import {MaterialModule} from "../material.module";
import {FormsModule} from "@angular/forms";
import {MatButtonModule} from "@angular/material/button";


@NgModule({
  declarations: [
    EventComponent
  ],
    imports: [
        CommonModule,
        EventRoutingModule,
        MaterialModule,
        FormsModule,
        MatButtonModule,
    ]
})
export class EventModule { }
