import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FederatedComponent } from './federated.component';

describe('FederatedComponent', () => {
  let component: FederatedComponent;
  let fixture: ComponentFixture<FederatedComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ FederatedComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(FederatedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
