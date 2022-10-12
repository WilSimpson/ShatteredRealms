import {TestBed} from '@angular/core/testing';

import {AuthGuard} from './auth.guard';
import {RouterTestingModule} from "@angular/router/testing";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {ActivatedRouteSnapshot, Router, RouterStateSnapshot} from "@angular/router";
import {AuthenticationService} from "../_services/authentication.service";
import {AlertService} from "../_services/alert.service";
import {Role} from "../models/role.model";

describe('AuthGuard', () => {
  let guard: AuthGuard;
  let router = {
    navigate: jasmine.createSpy('navigate'),
  };
  let alertService = {
    info: jasmine.createSpy('info'),
    error: jasmine.createSpy('error'),
  };
  let authenticationService: AuthenticationService;
  let routerStateSnapshot = jasmine.createSpyObj<RouterStateSnapshot>('RouterStateSnapshot', [], {url: 'https://test.com/'});
  let activatedRouteSnapshot = jasmine.createSpyObj<ActivatedRouteSnapshot>(
    'ActivatedRouterSnapshot',
    [],
    {data: {['role']: 'ADMIN'}},
  );

  const promise = Promise.resolve();

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [
        RouterTestingModule,
        HttpClientTestingModule,
      ],
      providers: [
        {provide: Router, useValue: router},
        {provide: AlertService, useValue: alertService},
        {provide: RouterStateSnapshot, useValue: routerStateSnapshot},
        {provide: ActivatedRouteSnapshot, useValue: activatedRouteSnapshot},
      ],
    }).compileComponents();
    guard = TestBed.inject(AuthGuard);
    authenticationService = TestBed.inject(AuthenticationService);
  });

  it('should be created', () => {
    expect(guard).toBeTruthy();
  });

  it('should succeed if there are no permissions required', () => {
    expect(guard['userHasPermission']([], [])).toBeTruthy();
    expect(guard['userHasPermission'](null, [])).toBeTruthy();
    expect(guard['userHasPermission'](null, null)).toBeTruthy();
  });

  it('should fail with no authenticated user', () => {
    spyOnProperty(authenticationService, 'currentUserValue').and.returnValue(undefined);
    expect(guard.canActivate(activatedRouteSnapshot, routerStateSnapshot)).toBeFalsy();
    expect(alertService.info).toHaveBeenCalledOnceWith('Please login to view this page');
    expect(router.navigate).toHaveBeenCalledWith(['/login'], {queryParams: {returnUrl: routerStateSnapshot.url}});
  });

  it('fail if they do not have the correct permissions', () => {
    spyOnProperty(authenticationService, 'currentUserValue').and.returnValue({
      email: "",
      firstName: "",
      id: 0,
      lastName: "",
      token: "",
      username: "",
      permissions: [],
      role:
    });

    router.navigate.and.returnValue(promise);
    expect(guard.canActivate(activatedRouteSnapshot, routerStateSnapshot)).toBeFalsy();
    expect(router.navigate).toHaveBeenCalledWith(['/']);
    promise.then(() => {
      expect(alertService.error).toHaveBeenCalledWith('ERROR: Access denied!');
    })
  });

  it('fail if they do not have the correct role', () => {
    spyOnProperty(authenticationService, 'currentUserValue').and.returnValue({
      email: "",
      firstName: "",
      id: 0,
      lastName: "",
      token: "",
      username: "",
      permissions: [],
      role: Role.USER
    });

    router.navigate.and.returnValue(promise);
    expect(guard.canActivate(activatedRouteSnapshot, routerStateSnapshot)).toBeFalsy();
    expect(router.navigate).toHaveBeenCalledWith(['/']);
    promise.then(() => {
      expect(alertService.error).toHaveBeenCalledWith('ERROR: Access denied!');
    })
  });

  it('succeed if they do not have the correct role and permissions', () => {
    spyOnProperty(authenticationService, 'currentUserValue').and.returnValue({
      email: "",
      firstName: "",
      id: 0,
      lastName: "",
      token: "",
      username: "",
      permissions: [Permission.TEST_PERMISSION],
      role: Role.ADMIN
    });

    router.navigate.and.returnValue(promise);
    expect(guard.canActivate(activatedRouteSnapshot, routerStateSnapshot)).toBeTruthy();
  });
});
