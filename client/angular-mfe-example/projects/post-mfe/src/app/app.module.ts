import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import {HttpClientModule} from "@angular/common/http";
import {APOLLO_NAMED_OPTIONS, NamedOptions} from "apollo-angular";
import {HttpLink} from "apollo-angular/http";
import {InMemoryCache, split} from "@apollo/client/core";
import { FileUploadTestComponent } from './file-upload-test/file-upload-test.component';
import {FormsModule} from "@angular/forms";
import {getMainDefinition} from "@apollo/client/utilities";
import {WebSocketLink} from "@apollo/client/link/ws";

@NgModule({
  declarations: [
    AppComponent,
    FileUploadTestComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    FormsModule
  ],
  providers: [
    {
      provide: APOLLO_NAMED_OPTIONS,
      useFactory(httpLink: HttpLink): NamedOptions {
        const http = httpLink.create({
          uri: 'http://localhost:8083/query',
        });


        const ws = new WebSocketLink({
          uri: `ws://localhost:8083/query`,
          options: {
            reconnect: true,
          }
        });

        const link = split(
          // split based on operation type
          ({query}) => {
            const data = getMainDefinition(query);
            return (
              data.kind === 'OperationDefinition' && data.operation === 'subscription'
            );
          },
          ws,
          http,
        );


        return {
          post: {
            cache: new InMemoryCache(),
            link: link
          }
        };
      },
      deps: [HttpLink],
    },
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
