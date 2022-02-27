import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { NotificationRoutingModule } from './notification-routing.module';
import { NotificationComponent } from './notification.component';
import {MatIconModule} from "@angular/material/icon";
import {MatButtonModule} from "@angular/material/button";

const EXPORTS = [
  NotificationComponent,
];

@NgModule({
  declarations: [
    ...EXPORTS,
  ],
  imports: [
    CommonModule,
    NotificationRoutingModule,
    MatIconModule,
    MatButtonModule
  ],
  exports: [...EXPORTS],
})
export class NotificationModule {
  static exports = EXPORTS;
}
