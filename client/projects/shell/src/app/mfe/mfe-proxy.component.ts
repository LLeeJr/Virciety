// src: https://github.com/manfredsteyer/module-federation-with-angular-dynamic-workflow-designer/

import {
  Component,
  ComponentFactoryResolver,
  Injector, Input,
  OnChanges,
  ViewChild,
  ViewContainerRef
} from '@angular/core';
import {loadRemoteModule} from "@angular-architects/module-federation";
import {MfeOptions} from "./mfe";

@Component({
  selector: 'mfe-proxy',
  template: `
        <ng-container #placeHolder></ng-container>
    `
})
export class MfeProxyComponent implements OnChanges {
  @ViewChild('placeHolder', { read: ViewContainerRef, static: true })
  viewContainer: ViewContainerRef;

  constructor(
    private injector: Injector,
    private cfr: ComponentFactoryResolver) { }

  @Input() options: MfeOptions;

  async ngOnChanges() {
    this.viewContainer.clear();

    const Component = await loadRemoteModule(this.options)
      .then(m => m[this.options.componentName]);

    const factory = this.cfr.resolveComponentFactory(Component);
    this.viewContainer.createComponent(factory, undefined, this.injector);
  }
}
