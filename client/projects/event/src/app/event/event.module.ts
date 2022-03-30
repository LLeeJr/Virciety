import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { EventRoutingModule } from './event-routing.module';
import {DialogReportCovidCaseComponent, DialogSubscribersComponent, EventComponent} from './event.component';
import {MaterialModule} from "../material.module";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import { CreateEventComponent } from '../create-event/create-event.component';
import { ContactDetailsComponent } from '../contact-details/contact-details.component';
import { SingleEventComponent } from '../single-event/single-event.component';


const EXPORTS = [
  SingleEventComponent,
  EventComponent
]

@NgModule({
  declarations: [
    ...EXPORTS,
    CreateEventComponent,
    DialogSubscribersComponent,
    ContactDetailsComponent,
    DialogReportCovidCaseComponent,
  ],
  imports: [
    CommonModule,
    EventRoutingModule,
    MaterialModule,
    FormsModule,
    ReactiveFormsModule,
    MaterialModule,
  ],
  exports: [...EXPORTS]
})
export class EventModule {
  static exports = EXPORTS
}
