<template>
  <div>
    <h2>Goal</h2>
    <p>
      Main goal of the WWN manager is to tackle the issue of identifying the hosts in SAN environment
      with the secondary goal being able to handle multitenant environments with potentially duplicate WWN.
    </p>

    <h2>Design</h2>
    <p>
      To identify the hostname the idea is to use SAN zoning, alias or WWN data. Every WWN in SAN may be associated
      with mutliple zones and preferably one alias. Zones and aliases follow naming convention that helps to identify
      the host the WWN belongs to. The naming convention can be different for different tenants or even not being followed
      for every zone/alias within the tenant. The tool uses a set of rules to identify the WWN and host using that WWN. In the end 
      the tool should generate a list of WWNs and hostname to be used in external tool for mapping WWNs to hosts.
    </p>

    <h3>WWN source</h3>
    <p>
      WWNs are imported to the tool in specific format containig following fields:
      <ul>
      <li>customer</li>
      <li>WWN</li>
      <li>zone</li>
      <li>alias</li>
      <li>existing hostname</li>
      <li>flag if the WNW is loaded by CSV from the source or detected</li>
      <li>WWN set type</li>
      </ul>

      WWN set type:
      <ul>
      <li>1 - WWN loaded from SAN configuration</li>
      <li>2 - WWN manually loaded into the source system</li>
      <li>3 - WWN automatically loaded into the source system</li>
      </ul>
      There be only 1 set 2 and set 3 type record for any given WWN and multiple set 1 records
    </p>

    <h3>Range Rules</h3>
    <p>
      Range rules are rules applied on global level. Their purpose is to separate host WWNs from WWNs used for array, backup or 
      other non-host purposes. They use regular expresion to match given WWN to one type of array, backup, host or other. The
      Range rules are applied first. Range rules are tested based on the defined order with first one matching being the rule to apply.
    </p>

    <h3>Host Rules</h3>
    <p>
      Host rules are applied only to Host WWNs and user regular expression to pluck the subset of the zone or alias name
      as the hostname for the given WWN. Zone and Alias base rules needs to contain match groups and specify which match 
      group denotes the host. Alternatively the whole WWN can be used to map WWN to a given hostname. The Host rules
      are applied second. Host rules can be either global or dedicated to a tenant. They are tested in defined order with
      customer rules being applied before global rules and first one matching being the one to apply.
    </p>
    
    <h3>Reconciliation Rules</h3>
    <p>
      These are special set of rules to handle duplicate WWNs accross multiple tenants or handle name missmatch between
      current hostnames and hostnames decoded using the Host rules. These rules are applied last. Additionaly there is 
      set of default Reconciliation rules that are always applied. Reconciliation rules can be either global or dedicated to a tenant.
      The order in wich they are applied is - default rules, customer defined rules and global rules. They are applied as long as
      the record still requires reconciliation. If the record requires reconciliation after all rules are applied, then manual 
      reconciliation needs to be done in the interface.
    </p>

    <h4>Default Reconciliation Rules</h4>
    <p>
      <ul>
      <li>Rule 1 - if wwn set is 3 and wwn belongs to mutliple customers the one owning the set 3 is selected as primary customer. All other customers with the same WWN are automatically marked as secondary. Additional if any hostname matches the set 3 hostname such records is marked as ignored for the export purpose.</li>
      <li>Rule 2 - if wwn set is 2 and wwn belongs to multiple customers, if any wwn set 1 hostname matches the set 2 hostname the customer owning set 2 is selected as primary, the set 2 is marked as ignored and all otherse are marked as secondary.</li>
      <li>Rule 3 - same scenario as Rule 2 but all hostnames are unique, then manual reconciliation is required.</li>
      <li>Rule 4 - if wwn is set 1 and belongs to mutliple customers where all are set 1 records, then manual reconciliation is needed.</li>
      </ul>
    </p>
  </div>
</template>