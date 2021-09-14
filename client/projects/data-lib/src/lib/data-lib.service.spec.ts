import { TestBed } from '@angular/core/testing';

import { DataLibService } from 'data-lib';

describe('DataLibService', () => {
  let service: DataLibService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(DataLibService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
