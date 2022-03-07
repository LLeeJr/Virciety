import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { NotificationRoutingModule } from './notification-routing.module';
import { NotificationComponent } from './notification.component';
import {MatIconModule} from "@angular/material/icon";
import {MatButtonModule} from "@angular/material/button";
import {MatCardModule} from "@angular/material/card";
import {MatBadgeModule} from "@angular/material/badge";

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
    MatButtonModule,
    MatCardModule,
    MatBadgeModule
  ],
  exports: [...EXPORTS],
})
export class NotificationModule {
  static exports = EXPORTS;
}
