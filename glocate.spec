Name:           glocate
Version:        %(cat .release-version 2>/dev/null || echo "0.1.2")
Release:        1%{?dist}
Summary:        Modern file search tool to replace locate
License:        MIT
URL:            https://github.com/Gosayram/go-locate
Source0:        %{name}-%{version}-linux-amd64.tar.gz
Source1:        %{name}-%{version}-linux-arm64.tar.gz

BuildArch:      x86_64 aarch64
BuildRequires:  systemd-rpm-macros

Provides:       %{name} = %{version}-%{release}
Provides:       locate = %{version}-%{release}
Obsoletes:      mlocate < %{version}-%{release}
Obsoletes:      findutils-locate < %{version}-%{release}

# Disable debug package generation
%global debug_package %{nil}

# Go-specific macros
%global goipath github.com/Gosayram/go-locate
%global commit  %{commit_hash}

%description
glocate is a modern, fast file search tool designed to replace the traditional
locate command. It provides enhanced search capabilities with better performance
and more intuitive command-line interface.

Key features:
- Fast file indexing and searching
- Modern command-line interface
- Configuration file support
- Backward compatibility with locate

%prep
# Extract the appropriate binary based on architecture
%ifarch x86_64
%setup -q -T -b 0 -n %{name}-%{version}-linux-amd64
%endif
%ifarch aarch64
%setup -q -T -b 1 -n %{name}-%{version}-linux-arm64
%endif

%build
# No build needed - using pre-compiled binaries

%install
# Install binary
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}

# Install configuration file
install -Dpm 0644 glocate.toml.example %{buildroot}%{_sysconfdir}/%{name}/glocate.toml

# Install documentation
install -Dpm 0644 README.md %{buildroot}%{_docdir}/%{name}/README.md
install -Dpm 0644 CHANGELOG.md %{buildroot}%{_docdir}/%{name}/CHANGELOG.md
install -Dpm 0644 LICENSE %{buildroot}%{_docdir}/%{name}/LICENSE

# Create symlink for backward compatibility
ln -sf %{name} %{buildroot}%{_bindir}/locate

%post
# Update locate database after installation
if [ $1 -eq 1 ]; then
    echo "Creating initial locate database..."
    %{_bindir}/%{name} --update-db >/dev/null 2>&1 || true
fi

%preun
# Remove symlink before uninstall
if [ $1 -eq 0 ]; then
    rm -f %{_bindir}/locate
fi

%files
%license %{_docdir}/%{name}/LICENSE
%doc %{_docdir}/%{name}/README.md
%doc %{_docdir}/%{name}/CHANGELOG.md
%dir %{_sysconfdir}/%{name}
%config(noreplace) %{_sysconfdir}/%{name}/glocate.toml
%{_bindir}/%{name}
%{_bindir}/locate

%changelog
* %(date '+%a %b %d %Y') Build System <abdurakhman.rakhmankulov@gmail.com> - %(cat .release-version 2>/dev/null || echo "0.1.2")-1
- Updated to use pre-compiled binaries following Go packaging best practices
- Added proper architecture support for x86_64 and aarch64
- Improved package metadata and dependencies
- Added backward compatibility symlink

* Sunday Jun 22 2025 Build System <abdurakhman.rakhmankulov@gmail.com> - 0.1.2-1
- Initial package for glocate
- Modern file search tool to replace locate
- Configuration file support
- Backward compatibility with locate command
