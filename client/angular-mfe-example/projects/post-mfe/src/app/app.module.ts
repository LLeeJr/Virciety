import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import {HttpClientModule} from "@angular/common/http";
import {APOLLO_NAMED_OPTIONS, NamedOptions} from "apollo-angular";
import {HttpLink} from "apollo-angular/http";
import {InMemoryCache} from "@apollo/client/core";
import { FileUploadTestComponent } from './file-upload-test/file-upload-test.component';
import {AngularMaterialModule} from "./angular-material.module";
import {MatToolbarModule} from "@angular/material/toolbar";
import {FormsModule} from "@angular/forms";

@NgModule({
  declarations: [
    AppComponent,
    FileUploadTestComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    AngularMaterialModule,
    MatToolbarModule,
    FormsModule
  ],
  providers: [
    {
      provide: APOLLO_NAMED_OPTIONS, // <-- Different from standard initialization
      useFactory(httpLink: HttpLink): NamedOptions {
        return {
          post: {
            cache: new InMemoryCache(),
            link: httpLink.create({
              uri: 'http://localhost:8083/query',
            }),
          },
        };
      },
      deps: [HttpLink],
    },
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
