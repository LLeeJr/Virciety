import {
  Component,
  Input,
  Injector,
  OnInit,
  ViewChild,
  ViewContainerRef,
  ɵcreateInjector, ComponentFactoryResolver, Output, EventEmitter,
} from "@angular/core";
import {loadRemoteModule} from "../utils/federation-utils";

@Component({
  selector: "app-federated-component",
  templateUrl: "./federated.component.html",
  styleUrls: ["./federated.component.scss"],
})
export class FederatedComponent implements OnInit {
  @ViewChild("federatedComponent", { read: ViewContainerRef })
  federatedComponent: ViewContainerRef;
  @Input() remoteEntry: string;
  @Input() remoteName: string;
  @Input() exposedModule: string;
  @Input() componentName: string;
  @Output() errorHandle = new EventEmitter<any>();

  constructor(private injector: Injector,
              private componentFactoryResolver: ComponentFactoryResolver) {}

  ngOnInit(): void {
    loadRemoteModule({
      remoteEntry: this.remoteEntry,
      remoteName: this.remoteName,
      exposedModule: this.exposedModule,
    }).then((federated) => {
      const component = federated[this.exposedModule].exports.find(
        (e: any) => e.ɵcmp?.exportAs[0] === this.componentName
      );

      const factory = this.componentFactoryResolver.resolveComponentFactory(component);

      const { instance } = this.federatedComponent.createComponent(
        factory, undefined, ɵcreateInjector(federated[this.exposedModule], this.injector)
      );
    }).catch(err => {
      this.errorHandle.emit({error: err, component: this.remoteName})
    });
  }
}
