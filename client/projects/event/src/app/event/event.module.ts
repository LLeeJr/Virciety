import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { EventRoutingModule } from './event-routing.module';
import {DialogSubscribersComponent, EventComponent} from './event.component';
import {MaterialModule} from "../material.module";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import { CreateEventComponent } from '../create-event/create-event.component';
import { ContactDetailsComponent } from '../contact-details/contact-details.component';


@NgModule({
  declarations: [
    EventComponent,
    CreateEventComponent,
    DialogSubscribersComponent,
    ContactDetailsComponent
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
