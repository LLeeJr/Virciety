import { TestBed } from '@angular/core/testing';

import { GQLService } from './g-q-l.service';

describe('GQLService', () => {
  let service: GQLService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(GQLService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
