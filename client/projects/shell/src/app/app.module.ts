import {APP_INITIALIZER, NgModule} from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import {HttpClientModule} from "@angular/common/http";
import {initializeKeycloak} from "./init/keycloak-init.factory";
import {KeycloakAngularModule, KeycloakService} from "keycloak-angular";
import {DatePipe} from "@angular/common";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {MaterialModule} from "./material.module";
import { FederatedComponent } from './federated/federated.component';
import { ProfileComponent } from './profile/profile.component';
import {ReactiveFormsModule} from "@angular/forms";
import {SinglePostComponent} from "./single-post/single-post.component";
import { SingleEventComponent } from './single-event/single-event.component';
import { PostComponent } from './post/post.component';
import { ChatComponent } from './chat/chat.component';
import { EventComponent } from './event/event.component';
import { OpenChatComponent } from './open-chat/open-chat.component';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import {environment} from "../environments/environment";

@NgModule({
  declarations: [
    AppComponent,
    FederatedComponent,
    ProfileComponent,
    SinglePostComponent,
    SingleEventComponent,
    PostComponent,
    ChatComponent,
    EventComponent,
    OpenChatComponent,
    PageNotFoundComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    KeycloakAngularModule,
    BrowserAnimationsModule,
    MaterialModule,
    ReactiveFormsModule,
  ],
  providers: [
    DatePipe,
    {
      provide: APP_INITIALIZER,
      useFactory: initializeKeycloak,
      multi: true,
      deps: [KeycloakService],
    },
    {
      provide: 'environment',
      useValue: environment
    }
  ],
  exports: [
    FederatedComponent
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
