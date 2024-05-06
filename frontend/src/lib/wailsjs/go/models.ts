export namespace mekano {
	
	export class BillingStatistics {
	    file: string;
	    debito: number;
	    credito: number;
	    base: number;
	
	    static createFrom(source: any = {}) {
	        return new BillingStatistics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file = source["file"];
	        this.debito = source["debito"];
	        this.credito = source["credito"];
	        this.base = source["base"];
	    }
	}
	export class PaymentStatistics {
	    archivo: string;
	    rango_rc: string;
	    bancolombia: number;
	    davivienda: number;
	    susuerte: number;
	    payu: number;
	    efectivo: number;
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new PaymentStatistics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.archivo = source["archivo"];
	        this.rango_rc = source["rango_rc"];
	        this.bancolombia = source["bancolombia"];
	        this.davivienda = source["davivienda"];
	        this.susuerte = source["susuerte"];
	        this.payu = source["payu"];
	        this.efectivo = source["efectivo"];
	        this.total = source["total"];
	    }
	}

}

