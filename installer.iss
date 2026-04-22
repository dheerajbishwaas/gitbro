[Setup]
AppName=GitBro
AppVersion=1.0
DefaultDirName={autopf}\gitbro
OutputBaseFilename=gitbro-setup
PrivilegesRequired=lowest
Compression=lzma
SolidCompression=yes

[Files]
Source: "gitbro.exe"; DestDir: "{app}"; Flags: ignoreversion

[Registry]
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "PATH"; ValueData: "{olddata};{app}"; Check: NeedsAddPath('{app}')

[Icons]
Name: "{group}\Uninstall GitBro"; Filename: "{uninstallexe}"

[UninstallDelete]
Type: filesandordirs; Name: "{app}"

[Code]
function NeedsAddPath(Param: string): boolean;
var
  OrigPath: string;
begin
  if not RegQueryStringValue(HKEY_CURRENT_USER, 'Environment', 'PATH', OrigPath)
  then begin
    Result := True;
    exit;
  end;
  Result := Pos(';' + Param + ';', ';' + OrigPath + ';') = 0;
end;